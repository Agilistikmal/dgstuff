export const apiUrl = "http://localhost:8080/api";

export const sendRequest = async (url, options = {}) => {
  const timeout = 10000;

  const controller = new AbortController();
  const timeoutId = setTimeout(() => {
    controller.abort();
    throw new Error("Request timed out");
  }, timeout);

  try {
    const response = await fetch(`${apiUrl}${url}`, {
      ...options,
      signal: controller.signal,
    });
    clearTimeout(timeoutId);
    return response.json();
  } catch (error) {
    clearTimeout(timeoutId);
    throw error;
  }
};
