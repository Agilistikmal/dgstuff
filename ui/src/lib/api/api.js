export const apiUrl = "http://localhost:8080/api";

export const sendRequest = async (url, options = {}) => {
  const timeout = 10000;

  return Promise.race([
    fetch(`${apiUrl}${url}`, {
      ...options,
      credentials: "include",
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((data) => {
            throw new Error(
              `${response.status} ${response.statusText}: ${data.message}`
            );
          });
        }
        return response.json();
      })
      .catch((error) => {
        throw error;
      }),

    new Promise((_, reject) => {
      setTimeout(() => {
        reject(new Error("Request timed out"));
      }, timeout);
    }),
  ]);
};
