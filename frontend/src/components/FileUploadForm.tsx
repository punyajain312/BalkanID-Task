import { useState, type DragEvent } from "react";
import { uploadFiles } from "../api/files";
import { useAuth } from "../context/AuthContext";
import toast from "react-hot-toast";

export default function FileUploadForm({
  onUploadSuccess,
}: {
  onUploadSuccess: () => void;
}) {
  const { token } = useAuth();
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setSelectedFiles(Array.from(e.target.files));
    }
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setSelectedFiles(Array.from(e.dataTransfer.files));
  };

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
  };

  const handleUpload = async () => {
    if (!token || selectedFiles.length === 0) return;
    setUploading(true);
    try {
      await uploadFiles(selectedFiles, token);
      toast.success("Upload successful ✅");
      setSelectedFiles([]); // clear file list
      onUploadSuccess(); // ✅ trigger FileList refresh
    } catch (err: any) {
      console.error("Upload error:", err.response || err.message);
      toast.error("Upload failed ❌");
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-4">
      {/* Drag & Drop area */}
      <div
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        className="w-full p-6 border-2 border-dashed rounded-lg text-center bg-gray-50 hover:bg-gray-100 cursor-pointer"
      >
        Drag & Drop files here
      </div>

      {/* File input */}
      <input type="file" multiple onChange={handleFileChange} className="w-full" />

      {/* Preview selected files */}
      {selectedFiles.length > 0 && (
        <ul className="text-sm text-gray-600">
          {selectedFiles.map((f, i) => (
            <li key={i}>{f.name}</li>
          ))}
        </ul>
      )}

      {/* Upload button */}
      <button
        onClick={handleUpload}
        disabled={uploading || selectedFiles.length === 0}
        className="w-full bg-blue-500 text-white p-2 rounded disabled:bg-gray-400"
      >
        {uploading ? "Uploading..." : "Upload"}
      </button>
    </div>
  );
}