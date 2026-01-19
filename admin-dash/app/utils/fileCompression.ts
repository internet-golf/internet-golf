import { packTar, createGzipEncoder, type TarEntry } from "modern-tar";

async function fileToUint8Array(file: File): Promise<Uint8Array> {
  const buffer = await file.arrayBuffer();
  return new Uint8Array(buffer);
}

export async function compressFilesToTarGz(files: File[]): Promise<Blob> {
  const entries: TarEntry[] = [];

  for (const file of files) {
    const path = file.webkitRelativePath || file.name;
    const data = await fileToUint8Array(file);
    entries.push({
      header: { name: path, size: data.length },
      body: data,
    });
  }

  entries.sort((a, b) => a.header.name.localeCompare(b.header.name));

  const tarData = await packTar(entries);

  const tarStream = new ReadableStream({
    start(controller) {
      controller.enqueue(tarData);
      controller.close();
    },
  });

  const gzipStream = tarStream.pipeThrough(createGzipEncoder());
  const response = new Response(gzipStream);
  const gzippedBuffer = await response.arrayBuffer();

  return new Blob([gzippedBuffer], { type: "application/gzip" });
}

export function getFilesFromUploadList(fileList: { originFileObj?: File }[] | undefined): File[] {
  if (!fileList) return [];
  return fileList.map((item) => item.originFileObj).filter((f): f is File => f !== undefined);
}
