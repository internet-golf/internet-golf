/**
 * README: although this is a generated file, allegedly it's perfectly safe to
 * edit, and that's how you can add things like custom base URLs and logic for
 * headers. the custom logic in this file has "MARK" comments
 */

import {
  clearGolfAuthTOken as clearGolfAuthToken,
  getGolfAuthToken,
  setGolfAuthToken,
} from "../session";

export type GolfFetcherExtraProps = {
  /**
   * You can add some extra props to your generated fetchers.
   *
   * Note: You need to re-gen after adding the first property to
   * have the `GolfFetcherExtraProps` injected in `GolfComponents.ts`
   **/
};

// MARK: base url
const baseUrl = "/_golf";

export type ErrorWrapper<TError> = TError | { status: "unknown"; payload: string };

export type GolfFetcherOptions<TBody, THeaders, TQueryParams, TPathParams> = {
  url: string;
  method: string;
  body?: TBody;
  headers?: THeaders;
  queryParams?: TQueryParams;
  pathParams?: TPathParams;
  signal?: AbortSignal;
} & GolfFetcherExtraProps;

const isBlob = (value: unknown): value is Blob => {
  if (typeof Blob !== "undefined" && value instanceof Blob) return true;
  return false;
};

const shouldUseFormData = (value: unknown): value is Record<string, unknown> => {
  if (!value || typeof value !== "object") return false;
  if (value instanceof FormData) return true;
  // If any value is Blob/File, treat it as multipart.
  return Object.values(value as Record<string, unknown>).some((v) => isBlob(v));
};

const objectToFormData = (value: Record<string, unknown>): FormData => {
  const formData = new FormData();
  for (const [key, v] of Object.entries(value)) {
    if (v === undefined || v === null) continue;
    if (isBlob(v)) {
      // If it's a File, keep its name; otherwise provide a stable filename.
      const fileName = (v as unknown as File).name || "upload";
      formData.append(key, v, fileName);
      continue;
    }
    if (typeof v === "string" || typeof v === "number" || typeof v === "boolean") {
      formData.append(key, String(v));
      continue;
    }
    // Fallback: serialize nested objects/arrays.
    formData.append(key, JSON.stringify(v));
  }
  return formData;
};

export async function golfFetch<
  TData,
  TError,
  TBody extends {} | FormData | undefined | null,
  THeaders extends {},
  TQueryParams extends {},
  TPathParams extends {},
>({
  url,
  method,
  body,
  headers,
  pathParams,
  queryParams,
  signal,
}: GolfFetcherOptions<TBody, THeaders, TQueryParams, TPathParams>): Promise<TData> {
  let error: ErrorWrapper<TError>;
  try {
    // MARK: auth token lookup
    const authToken = getGolfAuthToken();
    if (!authToken) {
      // hard navigate home. it would be nice to use react-router client-side
      // navigation, but it seems like even in the best case where we could
      // return something that could then turn into a call to redirect() or
      // navigate() in a loader or component, we would still have to check for
      // that in the loaders and components every time an api call is made. hard
      // navigation is probably fine.
      console.warn("no auth token found, redirecting home");
      window.location.href = "/";
    }
    const authHeader = { Authorization: `Bearer ${authToken}` };

    const requestHeaders: HeadersInit = {
      "Content-Type": "application/json",
      ...authHeader,
      ...headers,
    };

    const resolvedBody = (() => {
      if (!body) return undefined;
      if (body instanceof FormData) return body;
      if (shouldUseFormData(body)) return objectToFormData(body as Record<string, unknown>);
      return body;
    })();

    /**
     * As the fetch API is being used, when multipart/form-data is specified
     * the Content-Type header must be deleted so that the browser can set
     * the correct boundary.
     * https://developer.mozilla.org/en-US/docs/Web/API/FormData/Using_FormData_Objects#sending_files_using_a_formdata_object
     */
    // If we are sending FormData, the browser must set the Content-Type with boundary.
    if (resolvedBody instanceof FormData) {
      delete requestHeaders["Content-Type"];
    } else if (requestHeaders["Content-Type"]?.toLowerCase().includes("multipart/form-data")) {
      delete requestHeaders["Content-Type"];
    }

    const response = await window.fetch(`${baseUrl}${resolveUrl(url, queryParams, pathParams)}`, {
      signal,
      method: method.toUpperCase(),
      body: resolvedBody
        ? resolvedBody instanceof FormData
          ? resolvedBody
          : JSON.stringify(resolvedBody)
        : undefined,
      headers: requestHeaders,
    });
    if (!response.ok) {
      // MARK: redirect on 401
      if (response.status === 401) {
        // same comment as on auth token check. TODO: add a query parameter that
        // specifies the error and display it on the login page in red text
        console.warn("got 401 response, erasing token and redirecting home");
        clearGolfAuthToken();
        window.location.href = "/";
      }

      try {
        error = await response.json();
      } catch (e) {
        error = {
          status: "unknown" as const,
          payload: e instanceof Error ? `Unexpected error (${e.message})` : "Unexpected error",
        };
      }
    } else if (response.headers.get("content-type")?.includes("json")) {
      return await response.json();
    } else {
      // if it is not a json response, assume it is a blob and cast it to TData
      return (await response.blob()) as unknown as TData;
    }
  } catch (e) {
    const errorObject: Error = {
      name: "unknown" as const,
      message: e instanceof Error ? `Network error (${e.message})` : "Network error",
      stack: e as string,
    };
    throw errorObject;
  }
  throw error;
}

const resolveUrl = (
  url: string,
  queryParams: Record<string, string> = {},
  pathParams: Record<string, string> = {},
) => {
  let query = new URLSearchParams(queryParams).toString();
  if (query) query = `?${query}`;
  return url.replace(/\{\w*\}/g, (key) => pathParams[key.slice(1, -1)] ?? "") + query;
};
