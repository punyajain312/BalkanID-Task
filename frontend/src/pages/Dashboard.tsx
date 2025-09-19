import { useEffect, useState } from "react";
import { fetchFiles } from "../api/auth";
import { useAuth } from "../context/AuthContext";

export default function Dashboard() {
  const { token, logout } = useAuth();
  const [files, setFiles] = useState<any>(null);

  useEffect(() => {
    if (token) {
      fetchFiles(token)
        .then((res) => setFiles(res.data))
        .catch((err) => console.error("Error fetching files:", err));
    }
  }, [token]);

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold">Dashboard</h1>
      <button
        onClick={logout}
        className="bg-red-500 text-white px-3 py-1 rounded mt-2"
      >
        Logout
      </button>
      <pre className="mt-4 bg-gray-100 p-4 rounded">
        {JSON.stringify(files, null, 2)}
      </pre>
    </div>
  );
}