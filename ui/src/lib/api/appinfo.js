import { sendRequest } from "./api";

export class AppInfoApi {
  static basePath = "/appinfo";

  static async get() {
    const response = await sendRequest(this.basePath, {
      method: "GET",
    });
    return response;
  }
}
