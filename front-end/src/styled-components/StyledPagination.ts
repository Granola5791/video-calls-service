import { IconButton, styled } from "@mui/material";
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import ArrowBackIosNewIcon from '@mui/icons-material/ArrowBackIosNew';

const BasePaginationButton = styled(IconButton)({
    position: "absolute",
    top: "50%",
    transform: "translateY(-50%)",
    opacity: 0,
    transition: "opacity 0.3s ease-in-out",
    pointerEvents: "none",
    backgroundColor: "rgba(0, 0, 0, 0.5)",
});

export const StyledPrevPaginationButton = styled(BasePaginationButton)({
    right: "10px",
});

export const StyledNextPaginationButton = styled(BasePaginationButton)({
    left: "10px",
});

export const PaginationContainer = styled('div')({
    position: "absolute",
    top: 0,
    left: 0,
    height: "100%",
    width: "100%",
    '&:hover .meeting-pagination-btn': {
        opacity: 1,
        pointerEvents: "auto",
    },
});

export const StyledNextArrow = styled(ArrowForwardIosIcon)({
    fontSize: "2rem",
    color: "white",
});

export const StyledPrevArrow = styled(ArrowBackIosNewIcon)({
    fontSize: "2rem",
    color: "white",
});