import { useEffect, useState } from "react";
import { listFiles, deleteFile, generatePublicLink } from "../api/files";
import { useAuth } from "../context/AuthContext";
import toast from "react-hot-toast";

interface FileItem {
  id: string;
  filename: string;
  mime_type: string;
  size: number;
  created_at: string;
  hash: string;
}

export default function FileList({ refreshSignal }: { refreshSignal: number }) {
  const { token } = useAuth();
  const [files, setFiles] = useState<FileItem[]>([]);

  const loadFiles = async () => {
    if (!token) return;
    try {
      const res = await listFiles(token);
      setFiles(res.data.files || res.data);
    } catch (err) {
      toast.error("Failed to fetch files");
    }
  };

  const handleDelete = async (id: string) => {
    if (!token) return;
    try {
      await deleteFile(id, token);
      toast.success("File deleted");
      await loadFiles(); // ✅ refresh after delete
    } catch {
      toast.error("Delete failed");
    }
  };

  const handleShare = async (id: string) => {
    if (!token) return;
    try {
      const res = await generatePublicLink(id, token);
      const link = res.data.link;
      navigator.clipboard.writeText(link);
      toast.success("Public link copied!");
    } catch {
      toast.error("Failed to generate share link");
    }
  };

  useEffect(() => {
    loadFiles();
  }, [token, refreshSignal]); // ✅ reload after upload/delete

  return (
    <div className="mt-6">
      <h2 className="text-lg font-bold mb-2">Your Files</h2>
      {files.length === 0 ? (
        <p className="text-gray-500">No files uploaded yet.</p>
      ) : (
        <ul className="space-y-2">
          {files.map((file) => (
            <li
              key={file.id}
              className="p-2 border rounded flex justify-between items-center"
            >
              <div>
                <p className="font-medium">{file.filename}</p>
                <p className="text-sm text-gray-500">
                  {file.mime_type} · {file.size} bytes ·{" "}
                  {new Date(file.created_at).toLocaleString()}
                </p>
              </div>
              <div className="space-x-2">
                <button
                  onClick={() => handleShare(file.id)}
                  className="px-2 py-1 bg-yellow-500 text-white text-sm rounded"
                >
                  Share
                </button>
                <button
                  onClick={() => handleDelete(file.id)}
                  className="px-2 py-1 bg-red-500 text-white text-sm rounded"
                >
                  Delete
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}