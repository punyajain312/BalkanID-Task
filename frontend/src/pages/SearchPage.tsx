import { useState } from "react";
import SearchFilter from "../components/SearchFilter";

export default function SearchPage() {
  const [files, setFiles] = useState<any[]>([]);

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-2xl mx-auto bg-white p-6 rounded shadow">
        <SearchFilter onResults={setFiles} />
        <div className="mt-6">
          <h2 className="text-lg font-bold">Results</h2>
          {files.length === 0 ? (
            <p className="text-gray-500">No files found</p>
          ) : (
            <ul className="space-y-2 mt-2">
              {files.map((f, i) => (
                <li key={i} className="p-2 border rounded flex justify-between">
                  <span>{f.filename} ({f.mime_type})</span>
                  <span className="text-sm text-gray-500">{f.size} bytes</span>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}