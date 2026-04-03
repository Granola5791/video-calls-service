import { IconButton } from "@mui/material";
import DownloadIcon from '@mui/icons-material/Download';

interface DownloadButtonProps {
    fileName: string;
    putTextOnClick: () => string; // A function that returns the text to be downloaded when the button is clicked.
    className?: string;
}


const DownloadButton = ({ fileName, putTextOnClick, className }: DownloadButtonProps) => {

    const Download = () => {
        const text = putTextOnClick();
        const blob = new Blob([text], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);

        const link = document.createElement('a');
        link.href = url;
        link.download = fileName;

        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    };

    return (
        <IconButton className={className}>
            <DownloadIcon onClick={Download} />
        </IconButton>
    )
}

export default DownloadButton