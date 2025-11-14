import { sendRequest } from "./api";

export class StuffApi {
  static basePath = "/stuff";

  static async getAll(page = 1, limit = 10) {
    const response = await sendRequest(
      `${this.basePath}?page=${page}&limit=${limit}`,
      {
        method: "GET",
      }
    );
    return response;
  }

  static getFirstMedia(medias, mediaType = "image") {
    if (!medias) return null;
    return medias.find((media) => media.type === mediaType)?.url;
  }
}
