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

            <StyledPrevPaginationButton
                className="meeting-pagination-btn"
                color="inherit"
                onClick={onPrevClick}
                disabled={!hasPrev}
            >
                <StyledPrevArrow className="meeting-pagination-arrow" />
            </StyledPrevPaginationButton>

            <StyledNextPaginationButton
                className="meeting-pagination-btn"
                color="inherit"
                onClick={onNextClick}
                disabled={!hasNext}
            >
                <StyledNextArrow className="meeting-pagination-arrow" />
            </StyledNextPaginationButton>

        </PaginationContainer>
    )
}

export default MeetingPagination