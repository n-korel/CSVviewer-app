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
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <header className="mb-8 text-center">
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">
            CSV Viewer
          </h1>
          <p className="mt-2 text-lg text-gray-600">
            Загружайте и анализируйте большие CSV файлы
          </p>
        </header>

        <section className="mb-6 rounded-lg bg-white p-6 shadow">
          <FileUpload
            onUploadSuccess={handleUploadSuccess}
            setIsLoading={setIsLoading}
          />
        </section>

        {data && (
          <>
            <section className="mb-6 rounded-lg bg-white p-6 shadow">
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h2 className="text-2xl font-semibold text-gray-900">
                    Данные загружены
                  </h2>
                  <p className="mt-1 text-sm text-gray-600">
                    Всего записей:{" "}
                    <span className="font-medium">
                      {data.row_count?.toLocaleString()}
                    </span>
                  </p>
                </div>
                <button
                  onClick={handleClear}
                  className="inline-flex items-center justify-center rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                  type="button"
                >
                  Очистить данные
                </button>
              </div>

              <div className="mt-6">
                <SearchBar
                  searchQuery={searchQuery}
                  setSearchQuery={setSearchQuery}
                />
              </div>
            </section>

            <section className="rounded-lg bg-white p-6 shadow">
              <DataTable headers={data.headers} searchQuery={searchQuery} />
            </section>
          </>
        )}

        {isLoading && (
          <section className="rounded-lg bg-white p-12 text-center shadow">
            <div
              className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent"
              role="status"
            >
              <span className="sr-only">Загрузка...</span>
            </div>
            <p className="mt-4 text-gray-600">Загрузка файла...</p>
          </section>
        )}
      </main>
    </div>
  );
}

export default App;
