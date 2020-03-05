export type Client = {
  get<T = any>(path: string, options?: RequestOptions): Promise<T>;
  post<T = any>(path: string, body: any, options?: RequestOptions): Promise<T>;
};

export type ClientOptions = {
  endpoint: string;
  token: string;
};

export type RequestOptions = Partial<
  RequestInit & {
    path: string;
    body: any;
  }
>;

export const createClient = ({ endpoint, token }: ClientOptions): Client => {
  const request = ({ path, body, ...options }: RequestOptions) => {
    const params: RequestInit = {
      ...options,
    };

    if (body) {
      params.body = JSON.stringify(body);
    }

    if (token) {
      params.headers = {
        ...params.headers,
        Authorization: `Bearer ${token}`,
      };
    }

    return fetch(`${endpoint}${path}`, params)
      .then(response => {
        // We don't deal with malformed JSON at the moment.
        return response.json().then(json => ({
          response,
          json,
        }));
      })
      .then(({ response, json }) => {
        if (!response.ok) {
          return Promise.reject(json);
        }

        return json;
      });
  };

  return {
    get(path, options) {
      return request({
        ...options,
        method: 'GET',
        path,
      });
    },

    post(path, body, options) {
      return request({
        ...options,
        method: 'POST',
        path,
        body,
      });
    },
  };
};
