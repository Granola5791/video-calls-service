import { createTheme } from '@mui/material/styles';
import { heIL } from '@mui/material/locale';
import { Colors } from '../constants/general-contants';

declare module '@mui/material/styles' {
    interface PaletteColor {
        hover?: string;
    }

    interface SimplePaletteColorOptions {
        hover?: string;
    }

    interface TypeBackground {
        dark: string;
        primary: string;
    }
}

export const theme = createTheme(
    {
        direction: 'rtl',
        palette: Colors,
    },
    heIL,
);