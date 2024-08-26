import React from 'react';
import styled from "styled-components";
import MSlider from "@mui/material/Slider";
import Typography from "@mui/material/Typography";

// Styled Components
const Container = styled.div`
  width: 100%;
`;

const StyledTypography = styled(Typography)`
  && {
    color: rgba(0, 0, 0, 0.54);
  }
`;

const Render = ({ label, ...props }) => (
  <Container>
    <StyledTypography>{label}</StyledTypography>
    <MSlider {...props} />
  </Container>
);

export default Render;
