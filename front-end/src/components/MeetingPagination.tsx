import { PaginationContainer, StyledNextPaginationButton, StyledPrevPaginationButton } from "../styled-components/StyledPagination";
import { StyledNextArrow, StyledPrevArrow } from "../styled-components/StyledPagination";

interface MeetingPaginationProps {
    hasNext: boolean;
    hasPrev: boolean;
    onNextClick: () => void;
    onPrevClick: () => void;
}

const MeetingPagination = ({ hasNext, hasPrev, onNextClick, onPrevClick }: MeetingPaginationProps) => {

    return (
        <PaginationContainer>
            {hasPrev &&
                <StyledPrevPaginationButton className="meeting-pagination-btn" color="inherit" onClick={onPrevClick}>
                    <StyledPrevArrow />
                </StyledPrevPaginationButton>
            }
            {hasNext &&
                <StyledNextPaginationButton className="meeting-pagination-btn" color="inherit" onClick={onNextClick}>
                    <StyledNextArrow />
                </StyledNextPaginationButton>
            }
        </PaginationContainer>
    )
}

export default MeetingPagination