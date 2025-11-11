import { useState } from "react";

export default function FileUpload({ onUploadSuccess, setIsLoading }) {
  const [selectedFile, setSelectedFile] = useState(null);
  const [error, setError] = useState(null);

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;
    if (!file.name.endsWith(".csv")) {
      setError("Пожалуйста, выберите CSV-файл");
      setSelectedFile(null);
      return;
    }
    setSelectedFile(file);
    setError(null);
  };

  const handleUpload = async () => {
    if (!selectedFile) return;
    setIsLoading(true);

    const formData = new FormData();
    formData.append("file", selectedFile);

    try {
      const res = await fetch("/api/upload", {
        method: "POST",
        body: formData,
      });
      if (!res.ok) throw new Error("Ошибка при загрузке файла");
      const data = await res.json();
      onUploadSuccess(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center gap-6">
      <label className="flex items-center justify-center border-2 border-dashed border-gray-300 rounded-2xl p-8 text-center cursor-pointer hover:border-sky-400 transition-colors w-[500px] h-[100px]">
        <input
          type="file"
          accept=".csv,text/csv"
          onChange={handleFileChange}
          className="hidden"
        />
        {selectedFile ? (
          <p className="font-medium text-gray-800">{selectedFile.name}</p>
        ) : (
          <p className="text-gray-500">
            Перетащите CSV сюда или нажмите, чтобы выбрать
          </p>
        )}
      </label>

      {error && (
        <p className="text-rose-600 text-sm bg-rose-50 px-4 py-2 rounded-lg w-full text-center border border-rose-200">
          {error}
        </p>
      )}

      <button
        onClick={handleUpload}
        disabled={!selectedFile}
        className="px-6 py-3 bg-sky-600 text-white rounded-lg font-medium hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-sky-300 disabled:bg-gray-300 disabled:cursor-not-allowed transition"
      >
        Загрузить
      </button>
    </div>
  );
}
