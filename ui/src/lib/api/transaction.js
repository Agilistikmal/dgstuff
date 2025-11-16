import { sendRequest } from "./api";

export class TransactionApi {
  static basePath = "/transaction";

  static async create(payload) {
    const response = await sendRequest(this.basePath, {
      method: "POST",
      body: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
    });
    return response;
  }

  static async get(id, token = "") {
    const response = await sendRequest(`${this.basePath}/${id}`, {
      method: "GET",
      headers: {
        "X-Transaction-Token": token,
      },
    });
    return response;
  }
}
