import { useEffect, useState } from "react";
import { listFiles, deleteFile, generatePublicLink } from "../api/files";
import { useAuth } from "../context/AuthContext";
import { shareFilePublic } from "../api/public";
import toast from "react-hot-toast";

export interface FileItem {
  id: string;
  filename: string;
  mime_type: string;
  size: number;
  created_at: string;
  hash: string;
}

type Props = {
  files?: FileItem[];
  onDelete?: (id: string) => Promise<void>;
  onShare?: (id: string) => Promise<void>;
  refreshSignal?: number;
  limit?: number;
};

export default function FileList({
  files: controlledFiles,
  onDelete,
  onShare,
  refreshSignal = 0,
  limit,
}: Props) {
  const { token } = useAuth();
  const [files, setFiles] = useState<FileItem[]>([]);
  const [loading, setLoading] = useState(false);

  const isControlled = Array.isArray(controlledFiles);

  const loadFiles = async () => {
    if (!token) return;
    setLoading(true);
    try {
      const res = await listFiles(token);
      const fetched = res.data.files || [];
      setFiles(fetched);
    } catch (err) {
      console.error("listFiles error:", err);
      toast.error("Could not fetch files"); // ✅ only error if API fails
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!isControlled) loadFiles();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token, refreshSignal]);

  const handleDelete = async (id: string) => {
    if (!token) return;
    try {
      if (onDelete) {
        await onDelete(id);
      } else {
        await deleteFile(id, token);
        await loadFiles();
      }
      toast.success("File deleted");
    } catch (err) {
      console.error("delete error:", err);
      toast.error("Delete failed");
    }
  };

  const handleShare = async (id: string) => {
    if (!token) return;
    try {
      if (onShare) {
        await onShare(id);
      } else {
        const res = await generatePublicLink(id, token);
        const link = res.data.link;
        await navigator.clipboard.writeText(link);
        toast.success("Public link copied!");
      }
    } catch (err) {
      console.error("share error:", err);
      toast.error("Failed to generate share link");
    }
  };

  const usedFiles = isControlled ? controlledFiles! : files;
  const displayed = typeof limit === "number" ? usedFiles.slice(0, limit) : usedFiles;

  if (loading && !isControlled) {
    return <p className="text-gray-500">Loading files...</p>;
  }

  if (displayed.length === 0) {
    return <p className="text-gray-500">No files uploaded yet.</p>;
  }

  return (
    <ul className="space-y-2">
      {displayed.map((file) => (
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
              onClick={async () => {
                try {
                  if (!token) return;
                  const res = await shareFilePublic(file.id, token);
                  const link = window.location.origin + res.data.link;
                  await navigator.clipboard.writeText(link);
                  toast.success("Public link copied to clipboard!");
                } catch (err) {
                  console.error("Share public error:", err);
                  toast.error("Failed to share publicly");
                }
              }}
              className="px-2 py-1 bg-purple-500 text-white text-sm rounded hover:bg-purple-600"
            >
              Share Publicly
            </button>
            <button
              onClick={() => handleDelete(file.id)}
              className="px-2 py-1 bg-red-500 text-white text-sm rounded hover:bg-red-600"
            >
              Delete
            </button>
          </div>
        </li>
      ))}
    </ul>
  );
}