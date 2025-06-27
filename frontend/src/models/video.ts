interface Thumbnail {
    url: string;
    width: number;
    height: number;
}

interface Thumbnails {
    default?: Thumbnail;
    medium?: Thumbnail;
    [key: string]: Thumbnail | undefined;
}

interface ResourceId {
    kind?: string;
    videoId?: string;
}

export interface Video {
    id: number; // Assuming ID is optional and can be undefined
    title: string;
    publishedAt?: string;
    description?: string;
    videoOwnerChannelTitle?: string;
    thumbnails?: Thumbnails;
    resourceId?: ResourceId;
}