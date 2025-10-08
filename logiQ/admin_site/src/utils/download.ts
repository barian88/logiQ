export interface JsonDownloadOptions {
    space?: number;
}

export interface JsonDownloadHandle {
    url: string;
    revoke: () => void;
}

export const createJsonDownloadHandle = (data: unknown, options?: JsonDownloadOptions): JsonDownloadHandle => {
    const json = JSON.stringify(data, null, options?.space ?? 2);
    const blob = new Blob([json], {type: 'application/json'});
    const url = URL.createObjectURL(blob);

    return {
        url,
        revoke: () => URL.revokeObjectURL(url),
    };
};
