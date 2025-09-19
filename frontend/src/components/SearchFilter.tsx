import { useState } from "react";
import { searchFiles } from "../api/files";
import { useAuth } from "../context/AuthContext";
import toast from "react-hot-toast";

export default function SearchFilter({ onResults }: { onResults: (files: any[]) => void }) {
  const { token } = useAuth();
  const [filters, setFilters] = useState({
    filename: "",
    mime: "",
    uploader: "",
    size_min: "",
    size_max: "",
    date_from: "",
    date_to: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilters({ ...filters, [e.target.name]: e.target.value });
  };

  const handleSearch = async () => {
    if (!token) return;
    try {
      const res = await searchFiles(filters, token);
      onResults(res.data.files);
    } catch (err) {
      console.error(err);
      toast.error("Search failed ❌");
    }
  };

  return (
    <div className="space-y-2 border p-4 rounded bg-gray-50">
      <h2 className="text-lg font-bold">Search & Filters</h2>
      <input name="filename" placeholder="Filename" value={filters.filename} onChange={handleChange} className="w-full p-2 border rounded" />
      <input name="mime" placeholder="MIME Type" value={filters.mime} onChange={handleChange} className="w-full p-2 border rounded" />
      <input name="uploader" placeholder="Uploader Email" value={filters.uploader} onChange={handleChange} className="w-full p-2 border rounded" />
      <input name="size_min" placeholder="Min Size (bytes)" value={filters.size_min} onChange={handleChange} className="w-full p-2 border rounded" />
      <input name="size_max" placeholder="Max Size (bytes)" value={filters.size_max} onChange={handleChange} className="w-full p-2 border rounded" />
      <input name="date_from" type="date" value={filters.date_from} onChange={handleChange} className="w-full p-2 border rounded" />
      <input name="date_to" type="date" value={filters.date_to} onChange={handleChange} className="w-full p-2 border rounded" />
      <button onClick={handleSearch} className="w-full bg-blue-500 text-white p-2 rounded">Search</button>
    </div>
  );
}