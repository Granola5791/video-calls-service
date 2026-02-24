import styled from "@emotion/styled";
import VisibilityIcon from '@mui/icons-material/Visibility';
import MenuIcon from '@mui/icons-material/Menu';
import MicOffIcon from '@mui/icons-material/MicOff';
import MicIcon from '@mui/icons-material/Mic';

export const HamburgerMenuIcon = styled(MenuIcon)({
    fontSize: "150%",
});

export const WatchingIcon = styled(VisibilityIcon)({
    color: "red",
    position: "absolute",
    left: "10%",
    fontSize: "250%",
});

export const MutedIcon = styled(MicOffIcon)({
    color: "gray",
    fontSize: "175%",
});

export const UnmutedIcon = styled(MicIcon)({
    color: "white",
    fontSize: "175%",
});

export const SmallMutedIcon = styled(MicOffIcon)({
    color: "gray",
    fontSize: "100%",
});

export const SmallUnmutedIcon = styled(MicIcon)({
    color: "white",
    fontSize: "100%",
});