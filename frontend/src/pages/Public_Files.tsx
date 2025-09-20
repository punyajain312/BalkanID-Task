import { useEffect, useState } from "react";
import { listPublicFiles, downloadPublicFile } from "../api/public";
import toast from "react-hot-toast";

interface PublicFile {
  id: string;
  filename: string;
  mime_type: string;
  size: number;
  created_at: string;
  shared_by: string;
  public_url: string;
}

export default function PublicFilesPage() {
  const [files, setFiles] = useState<PublicFile[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await listPublicFiles();
        setFiles(res.data.files || []);
      } catch (err) {
        toast.error("Failed to load public files");
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) return <p className="text-gray-500">Loading public files...</p>;

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-4xl mx-auto bg-white p-4 rounded shadow">
        <h1 className="text-xl font-bold mb-4">Publicly Shared Files</h1>
        {files.length === 0 ? (
          <p className="text-gray-500">No public files available.</p>
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
                <button
                  onClick={() => downloadPublicFile(file.public_url)}
                  className="px-2 py-1 bg-blue-500 text-white text-sm rounded hover:bg-blue-600"
                >
                  Download
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}