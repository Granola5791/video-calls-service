import styled from "@emotion/styled";
import { TextField } from "@mui/material";

export const StyledTextField = styled(TextField)({
    width: "300px",
    textAlign: "center",
    '& .MuiInputLabel-root': {
        // Position the label on the right instead of the left
        right: '16px',
        left: 'auto',
        // Ensure it shrinks toward the right side
        transformOrigin: 'top right',
    },
    '& .MuiInputLabel-shrink': {
        transform: 'translate(0, -8px) scale(0.75) !important', // keeps label vertically centered in the gap
        right: '15px', // aligns label with notched outline gap
    },
    '& .MuiOutlinedInput-notchedOutline': {
        textAlign: 'right', // moves the legend/gap to the right side
    },
    '& .MuiFilledInput-root': {
        // Fixes gray gaps in the background of the text field
        backgroundColor: '#e3f2fd',
        '&.Mui-focused': {
            backgroundColor: '#e3f2fd',
        },
    },
});