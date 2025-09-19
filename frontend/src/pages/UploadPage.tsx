import FileUploadForm from "../components/FileUploadForm";
import FileList from "../components/FileList";

export default function UploadPage() {
  return (
    <div className="min-h-screen flex flex-col items-center bg-gray-100 p-6">
      <div className="w-full max-w-lg bg-white p-6 rounded shadow">
        <h1 className="text-2xl font-bold mb-4">Upload Files</h1>
        <FileUploadForm />
        <FileList />
      </div>
    </div>
  );
}