import axios from "axios";

export const axiosInstance = axios.create({
  baseURL: process.env.STOQR_API_URL ? process.env.STOQR_API_URL : "http://localhost:8080/",
  timeout: 1000,
});
