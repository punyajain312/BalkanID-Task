import axios from "axios";

const API_URL = "http://localhost:8080";

export const signup = (name: string, email: string, password: string) =>
  axios.post(`${API_URL}/signup`, { name, email, password });

export const login = (email: string, password: string) =>
  axios.post(`${API_URL}/login`, { email, password });

export const fetchFiles = (token: string) =>
  axios.get(`${API_URL}/files`, {
    headers: { Authorization: `Bearer ${token}` },
  });