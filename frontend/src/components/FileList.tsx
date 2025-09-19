import { useEffect, useState } from "react";
import { listFiles } from "../api/files";
import { useAuth } from "../context/AuthContext";

export default function FileList() {
  const { token } = useAuth();
  const [files, setFiles] = useState<any[]>([]);

  useEffect(() => {
    if (token) {
      listFiles(token)
        .then((res) => setFiles(res.data.files || []))
        .catch((err) => console.error("Error fetching files:", err));
    }
  }, [token]);

  return (
    <div className="mt-6">
      <h2 className="text-lg font-bold mb-2">Your Files</h2>
      {files.length === 0 ? (
        <p className="text-gray-500">No files uploaded yet.</p>
      ) : (
        <ul className="space-y-2">
          {files.map((file, idx) => (
            <li key={idx} className="p-2 border rounded flex justify-between">
              <span>{file.filename}</span>
              <span className="text-sm text-gray-500">{file.size} bytes</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}