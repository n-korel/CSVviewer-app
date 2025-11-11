import { useState } from "react";
import FileUpload from "./components/FileUpload";
import DataTable from "./components/DataTable";
import SearchBar from "./components/SearchBar";

function App() {
  const [data, setData] = useState(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleUploadSuccess = (uploadedData) => {
    setData(uploadedData);
    setSearchQuery("");
  };

  const handleClear = async () => {
    try {
      await fetch("/api/clear", { method: "DELETE" });
      setData(null);
      setSearchQuery("");
    } catch (error) {
      console.error("Error clearing data:", error);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <header className="mb-10 text-center">
          <h1 className="text-5xl font-bold text-gray-800 mb-2">CSV Viewer</h1>
          <p className="text-gray-600 text-lg">
            Загружайте и анализируйте большие CSV файлы
          </p>
        </header>

        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <FileUpload
            onUploadSuccess={handleUploadSuccess}
            setIsLoading={setIsLoading}
          />
        </div>

        {data && (
          <>
            <div className="bg-white rounded-lg shadow-md p-6 mb-6">
              <div className="flex justify-between items-center mb-4">
                <div>
                  <h2 className="text-2xl font-semibold text-gray-800">
                    Данные загружены
                  </h2>
                  <p className="text-gray-600 mt-1">
                    Всего записей: {data.row_count?.toLocaleString()}
                  </p>
                </div>
                <button
                  onClick={handleClear}
                  className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
                >
                  Очистить данные
                </button>
              </div>

              <SearchBar
                searchQuery={searchQuery}
                setSearchQuery={setSearchQuery}
              />
            </div>

            <div className="bg-white rounded-lg shadow-md p-6">
              <DataTable headers={data.headers} searchQuery={searchQuery} />
            </div>
          </>
        )}

        {isLoading && (
          <div className="bg-white rounded-lg shadow-md p-12 text-center">
            <p className="text-gray-600">Загрузка файла...</p>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
