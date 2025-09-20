import { Link, Routes, Route, useNavigate } from "react-router-dom";
import FileUploadForm from "../components/FileUploadForm";
import FileList, { type FileItem } from "../components/FileList";
import SearchFilter from "../components/SearchFilter";
import { useState, useEffect } from "react";
import { listFiles, deleteFile, generatePublicLink } from "../api/files";
import { useAuth } from "../context/AuthContext";
import toast from "react-hot-toast";

export default function Dashboard() {
  const { token } = useAuth();
  const navigate = useNavigate();

  const [files, setFiles] = useState<FileItem[]>([]);
  const [refresh, setRefresh] = useState(0);
  const [searchResults, setSearchResults] = useState<any[]>([]);

  const loadFiles = async () => {
    if (!token) return;
    try {
      const res = await listFiles(token);
      setFiles(res.data.files || res.data);
    } catch (err) {
      toast.error("Failed to load files");
    }
  };

  useEffect(() => {
    loadFiles();
  }, [token, refresh]);

  const handleDelete = async (id: string) => {
    if (!token) return;
    try {
      await deleteFile(id, token);
      setFiles((prev) => prev.filter((f) => f.id !== id));
      toast.success("File deleted");
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

  const onUploadSuccess = () => {
    setRefresh((r) => r + 1); // refresh files
    navigate("/dashboard"); // ✅ redirect to dashboard home
  };

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-5xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold mb-6">Dashboard</h1>

        {/* Navigation */}
        <div className="flex space-x-4 mb-6">
          <Link to="/dashboard/upload" className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
            Upload File
          </Link>
          <Link to="/dashboard/search" className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600">
            Search
          </Link>
          <Link to="/dashboard/list" className="px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600">
            My Files
          </Link>
          <Link to="/dashboard/public" className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600">
            Public Files
          </Link>
        </div>

        <div className="bg-white p-4 rounded shadow">
          <Routes>
            {/* Dashboard Home (Recent Files preview) */}
            <Route
              path="/"
              element={
                <div>
                  <h2 className="text-lg font-bold mb-4">Recent Files</h2>
                  <FileList files={files} onDelete={handleDelete} onShare={handleShare} limit={5} />
                </div>
              }
            />

            {/* Upload Page */}
            <Route
              path="upload"
              element={
                <div>
                  <h2 className="text-lg font-bold mb-4">Upload Files</h2>
                  <FileUploadForm onUploadSuccess={onUploadSuccess} />
                </div>
              }
            />

            {/* Search Page */}
            <Route
              path="search"
              element={
                <div>
                  <h2 className="text-lg font-bold mb-4">Search Files</h2>
                  <SearchFilter onResults={setSearchResults} />
                  {searchResults.length > 0 && (
                    <div className="mt-4">
                      <h3 className="font-bold">Search Results</h3>
                      <ul className="mt-2 space-y-2">
                        {searchResults.map((f, i) => (
                          <li key={i} className="p-2 border rounded flex justify-between items-center">
                            {f.filename} ({f.mime_type}) - {f.size} bytes
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              }
            />

            {/* Full File List */}
            <Route
              path="list"
              element={
                <div>
                  <h2 className="text-lg font-bold mb-4">My Files</h2>
                  <FileList files={files} onDelete={handleDelete} onShare={handleShare} />
                </div>
              }
            />
          </Routes>
        </div>
      </div>
    </div>
  );
}