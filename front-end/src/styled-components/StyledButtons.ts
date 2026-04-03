import { Button, IconButton, styled } from "@mui/material";

export const LongButton = styled(Button)({
    width: "300px",
    height: "50px",
    fontSize: "20px",
    fontWeight: "bold",
    borderRadius: "20px",
    border: `2px solid lightgray`,
    color: "gray",
});

export const LongButtonFilled = styled(LongButton)(({ theme }) => ({
    color: "white",
    backgroundColor: theme.palette.primary.main,
    border: `2px solid ${theme.palette.primary.hover}`,
    "&:hover": {
        backgroundColor: theme.palette.primary.hover,
    },
}));

export const MediumButton = styled(Button)({
    width: "200px",
    height: "40px",
    fontSize: "16px",
    fontWeight: "bold",
    borderRadius: "20px",
    border: `2px solid lightgray`,
    color: "gray",
});

export const SmallButton = styled(Button)({
    width: "100px",
    height: "35px",
    fontSize: "12px",
    fontWeight: "bold",
    borderRadius: "20px",
    border: `2px solid lightgray`,
    color: "gray",
});

export const AdaptiveButton = styled(Button)({
    width: 'fit-content',
    height: '35px',
    fontSize: '15px',
    fontWeight: 'bold',
    borderRadius: '20px',
    border: `2px solid lightgray`,
    color: 'gray',
    padding: '15px 20px',
})

export const AdaptiveButtonFilled = styled(AdaptiveButton)(({ theme }) => ({
    color: "white",
    backgroundColor: theme.palette.primary.main,
    border: `2px solid ${theme.palette.primary.hover}`,
    "&:hover": {
        backgroundColor: theme.palette.primary.hover,
    },
}));

export const SmallButtonFilled = styled(SmallButton)(({ theme }) => ({
    color: "white",
    backgroundColor: theme.palette.primary.main,
    border: `2px solid ${theme.palette.primary.hover}`,
    "&:hover": {
        backgroundColor: theme.palette.primary.hover,
    },
}));

export const MediumButtonFilled = styled(MediumButton)(({ theme }) => ({
    color: "white",
    backgroundColor: theme.palette.primary.main,
    border: `2px solid ${theme.palette.primary.hover}`,
    "&:hover": {
        backgroundColor: theme.palette.primary.hover,
    },
}));

export const LeaveMeetingButton = styled(MediumButton)(({ theme }) => ({
    color: "white",
    backgroundColor: theme.palette.error.main,
    border: `2px solid ${theme.palette.error.hover}`,
    "&:hover": {
        backgroundColor: theme.palette.error.hover,
    },
}));

export const HamburgerMenuButton = styled(IconButton)({
    position: "fixed",
    backgroundColor: "transparent",
    border: "none",
    left: "4%",
    top: "4%",
});

export const MuteButton = styled(IconButton)({
    position: "fixed",
    right: "5%",
});

export const SmallMuteButton = styled(IconButton)({
    position: "absolute",
    left: "3%",
    top: "5%",
    padding: "0",
});