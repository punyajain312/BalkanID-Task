import axios from "axios";

const API_URL = "http://localhost:8080";

export const uploadFiles = (files: File[], token: string) => {
  const formData = new FormData();
  files.forEach((file) => formData.append("files", file));

  return axios.post(`${API_URL}/upload`, formData, {
    headers: { Authorization: `Bearer ${token}` },
  });
};

export const listFiles = (token: string) =>
  axios.get(`${API_URL}/files`, {
    headers: { Authorization: `Bearer ${token}` },
  });

export const searchFiles = (params: Record<string, string>, token: string) =>
  axios.get(`${API_URL}/search`, {
    headers: { Authorization: `Bearer ${token}` },
    params,
  });

export const deleteFile = (id: string, token: string) =>
  axios.delete(`${API_URL}/delete?id=${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });

export const generatePublicLink = (id: string, token: string) =>
  axios.post(`${API_URL}/share`, { id }, {
    headers: { Authorization: `Bearer ${token}` },
  });