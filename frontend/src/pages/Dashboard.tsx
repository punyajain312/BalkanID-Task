import { useState } from "react";
import FileUploadForm from "../components/FileUploadForm";
import FileList from "../components/FileList";
import SearchFilter from "../components/SearchFilter";

export default function Dashboard() {
  const [refresh, setRefresh] = useState(0);
  const [activeSection, setActiveSection] = useState<"list" | "upload" | "search">("list");
  const [searchResults, setSearchResults] = useState<any[]>([]);

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-4xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold mb-4">Dashboard</h1>

        {/* Navigation */}
        <div className="flex space-x-4">
          <button
            className={`px-4 py-2 rounded ${activeSection === "upload" ? "bg-blue-500 text-white" : "bg-gray-200"}`}
            onClick={() => setActiveSection("upload")}
          >
            Upload File
          </button>
          <button
            className={`px-4 py-2 rounded ${activeSection === "search" ? "bg-blue-500 text-white" : "bg-gray-200"}`}
            onClick={() => setActiveSection("search")}
          >
            Search
          </button>
          <button
            className={`px-4 py-2 rounded ${activeSection === "list" ? "bg-blue-500 text-white" : "bg-gray-200"}`}
            onClick={() => setActiveSection("list")}
          >
            My Files
          </button>
        </div>

        {/* Sections */}
        {activeSection === "upload" && (
          <div className="bg-white p-4 rounded shadow">
            <h2 className="text-lg font-bold mb-2">Upload Files</h2>
            <FileUploadForm
              onUploadSuccess={() => {
                setRefresh((r) => r + 1);   
                setActiveSection("list"); 
              }}
            />
          </div>
        )}

        {activeSection === "search" && (
          <div className="bg-white p-4 rounded shadow">
            <SearchFilter onResults={setSearchResults} />
            {searchResults.length > 0 && (
              <div className="mt-4">
                <h2 className="font-bold">Search Results</h2>
                <ul className="mt-2 space-y-2">
                  {searchResults.map((f, i) => (
                    <li
                      key={i}
                      className="p-2 border rounded flex justify-between items-center"
                    >
                      <span>
                        {f.filename} ({f.mime_type}) - {f.size} bytes
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        )}

        {activeSection === "list" && (
          <div className="bg-white p-4 rounded shadow">
            <FileList refreshSignal={refresh} />
          </div>
        )}
      </div>
    </div>
  );
}