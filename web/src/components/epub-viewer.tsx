import { useState } from "react";
import { ReactReader } from "react-reader";

interface EpubViewerProps {
	fileUrl: string;
}

const EpubViewer = ({ fileUrl }: EpubViewerProps) => {
	const [location, setLocation] = useState<string | number>(0);

	return (
		<div style={{ height: "100vh" }}>
			<ReactReader
				url={fileUrl}
				location={location}
				locationChanged={(epubcfi: string) => setLocation(epubcfi)}
			/>
		</div>
	);
};

export default EpubViewer;
