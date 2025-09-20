import axios from "axios";
const API_URL = "http://localhost:8080";

// Share a file publicly
export const shareFilePublic = (fileId: string, token: string) =>
  axios.post(`${API_URL}/share?id=${fileId}`, null, {
    headers: { Authorization: `Bearer ${token}` },
  });

// List all public files
export const listPublicFiles = () =>
  axios.get(`${API_URL}/public/list`);

// Download public file by token
export const downloadPublicFile = (token: string) =>
  window.open(`${API_URL}/public/file?token=${token}`, "_blank");