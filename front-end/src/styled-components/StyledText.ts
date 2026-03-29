import { keyframes } from "@emotion/react";
import { styled } from "@mui/material";

const dotAnimation = keyframes`
  0%, 100% { content: "."; }
  33% { content: ".."; }
  66% { content: "..."; }
`;

export const LoadingText = styled('span')({
  '&::after': {
    display: 'inline-block',
    content: '"."',
    animation: `${dotAnimation} 1.5s steps(1, end) infinite`,
    width: '1em',
  },
});

export const StyledTitle = styled('h1')(({ theme }) => ({
    fontSize: '5rem',
    fontWeight: 'bold',
    color: theme.palette.primary.main,
    textAlign: 'center',
    margin: '0px',
}));

export const BoldText = styled('p')({
    fontWeight: 'bold',
});

export const StyledH1 = styled('h1')(({ theme }) => ({
    color: theme.palette.primary.main,
}));

export const StyledP = styled('p')(({ theme }) => ({
    color: theme.palette.primary.main,
    margin: '0px',
}));

export const MeetingID = styled('p')({
    position: 'absolute',
    left: '10px',
    color: 'gray',
    width: '30%',
    fontSize: '0.9rem',
});

export const NameTag = styled('p')({
    position: 'absolute',
    right: '0px',
    top: '0px',
    color: 'white',
    width: 'fit-content',
    height: 'fit-content',
    fontSize: '0.9rem',
    backgroundColor: 'black',
    margin: '0px',
    padding: '1%',
});

export const ClickableP = styled('p')(({ theme }) => ({
    cursor: 'pointer',
    color: theme.palette.primary.main,
    '&:hover': {
        color: theme.palette.primary.hover,
    }
}));