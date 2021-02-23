import axios from "axios";

export const axiosInstance = axios.create({
  baseURL: (window as any).STOQR_API_URL,
  timeout: 1000,
});
