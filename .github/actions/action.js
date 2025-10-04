const core = require('@actions/core');
const tc = require('@actions/tool-cache');

async function setup() {
    // Get version of tool to be installed
    const version = core.getInput('version');

    const os = process.env.RUNNER_OS;
    const exe = os === 'Windows' ? 'golf-windows-amd64.exe' : 'golf-linux-amd64';

    // Download the specific version of the tool, e.g. as a tarball
    const pathToCLI = await tc.downloadTool(
        `https://github.com/toBeOfUse/internet-golf/releases/download/${version}/${exe}`
    );

    // Expose the tool by adding it to the PATH
    core.addPath(pathToCLI)
}

module.exports = setup
