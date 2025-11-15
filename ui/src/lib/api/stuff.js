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

  static getFirstMedia(medias, mediaType = null) {
    if (!medias) return null;
    if (mediaType) {
      return medias.find((media) => media.type === mediaType);
    }
    return medias[0];
  }

  static async getBySlug(slug) {
    const response = await sendRequest(
      `${this.basePath}/${slug}`,
      {
        method: "GET",
      }
    );
    return response;
  }
}
