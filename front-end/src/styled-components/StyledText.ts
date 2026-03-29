import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";
import { Colors } from "../constants/general-contants.ts";

const dotAnimation = keyframes`
  0%, 100% { content: "."; }
  33% { content: ".."; }
  66% { content: "..."; }
`;

export const LoadingText = styled.span({
  '&::after': {
    display: 'inline-block',
    content: '"."',
    animation: `${dotAnimation} 1.5s steps(1, end) infinite`,
    width: '1em',
    textAlign: 'right',
  },
});

export const StyledTitle = styled.h1({
    fontSize: '5rem',
    fontWeight: 'bold',
    color: Colors.primary,
    textAlign: 'center',
    margin: '0px',
});

export const BoldText = styled.p({
    fontWeight: 'bold',
});

export const StyledH1 = styled.h1({
    color: Colors.primary,
});

export const StyledP = styled.p({
    color: Colors.primary,
    margin: '0px',
});

export const MeetingID = styled.p({
    position: 'absolute',
    right: '10px',
    color: 'gray',
    width: '30%',
    fontSize: '0.9rem',
});

export const NameTag = styled.p({
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

export const ClickableP = styled.p({
    cursor: 'pointer',
    color: Colors.primary,
    '&:hover': {
        color: Colors.primaryHover,
    }
})