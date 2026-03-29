import { IconButton, InputAdornment } from "@mui/material";
import { StyledTextField } from "../styled-components/StyledTextFields"
import VisibilityIcon from '@mui/icons-material/Visibility';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import { useState, type ChangeEvent } from "react";

interface PasswordFieldProps {
    label: string;
    variant?: "filled" | "outlined" | "standard";
    onChange: (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
}

const PasswordField = ({ label, variant, onChange }: PasswordFieldProps) => {
    const [showPassword, setShowPassword] = useState(false);

    const toggleShowPassword = () => {
        setShowPassword((prev) => !prev);
    }

    return (
        <StyledTextField
            label={label}
            type={showPassword ? 'text' : 'password'}
            variant={variant}
            onChange={onChange}
            slotProps={{
                input: {
                    endAdornment: (
                        <InputAdornment position="end">
                            <IconButton
                                onClick={toggleShowPassword}
                                onMouseDown={(e) => e.preventDefault()}
                                edge="end"
                            >
                                {showPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
                            </IconButton>
                        </InputAdornment>
                    )
                }
            }}
        />
    )
}

export default PasswordField