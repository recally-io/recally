// Core viewer
import { SpecialZoomLevel, Viewer, Worker } from "@react-pdf-viewer/core";
import { defaultLayoutPlugin } from "@react-pdf-viewer/default-layout";
import { pageNavigationPlugin } from "@react-pdf-viewer/page-navigation";
import { searchPlugin } from "@react-pdf-viewer/search";
import { thumbnailPlugin } from "@react-pdf-viewer/thumbnail";
import { zoomPlugin } from "@react-pdf-viewer/zoom";

// Import styles
import "@react-pdf-viewer/core/lib/styles/index.css";
import "@react-pdf-viewer/default-layout/lib/styles/index.css";
import "@react-pdf-viewer/page-navigation/lib/styles/index.css";
import "@react-pdf-viewer/search/lib/styles/index.css";
import "@react-pdf-viewer/thumbnail/lib/styles/index.css";
import "@react-pdf-viewer/zoom/lib/styles/index.css";

export default function PdfViewer({ fileUrl }: { fileUrl: string }) {
	// Initialize plugins
	const defaultLayoutPluginInstance = defaultLayoutPlugin();
	const zoomPluginInstance = zoomPlugin();
	const pageNavigationPluginInstance = pageNavigationPlugin();
	const searchPluginInstance = searchPlugin();
	const thumbnailPluginInstance = thumbnailPlugin();

	return (
		<div className="h-[calc(100vh-4rem)] w-full border rounded-lg shadow-sm overflow-hidden">
			<Worker workerUrl="https://unpkg.com/pdfjs-dist@3.4.120/build/pdf.worker.min.js">
				<Viewer
					fileUrl={fileUrl}
					defaultScale={SpecialZoomLevel.PageFit}
					plugins={[
						defaultLayoutPluginInstance,
						zoomPluginInstance,
						pageNavigationPluginInstance,
						searchPluginInstance,
						thumbnailPluginInstance,
					]}
				/>
			</Worker>
		</div>
	);
}
