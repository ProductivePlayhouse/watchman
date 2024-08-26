// PPH ADDED FILE
import React from "react";
import { makeStyles } from "@mui/styles";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Typography from '@mui/material/Typography';
import Paper from "@mui/material/Paper";

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
});

export const SectoralSanctions = ({ data }) =>
{
  const classes = useStyles();

  if (!data) return null;

  return (
    <TableContainer component={Paper}>
      <Typography variant="h5">Sectoral Sanctions</Typography>
      <Table className={classes.table} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Entity ID</TableCell>
            <TableCell>Name</TableCell>
            <TableCell>Type</TableCell>
            <TableCell>Programs</TableCell>
            <TableCell>Match</TableCell>
            <TableCell>Addresses</TableCell>
            <TableCell>Alternate Names</TableCell>
            <TableCell>IDs</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data.map((s) => (
            <TableRow key={s.EntityID}>
              <TableCell>{s.EntityID}</TableCell>
              <TableCell>{s.Name}</TableCell>
              <TableCell>{s.Type}</TableCell>
              <TableCell>
                {s.Programs.map((p, i) => (
                  <span key={i}>
                    {p}
                    {i < s.Programs.length - 1 && ", "}
                  </span>
                ))}
              </TableCell>
              <TableCell>{s.Match}</TableCell>
              <TableCell>
                {s.Addresses.map((p, i) => (
                  <span key={i}>
                    {p}
                    {i < s.Addresses.length - 1 && ", "}
                  </span>
                ))}
              </TableCell>
              <TableCell>
                {s.AlternateNames && s.AlternateNames.map((p, i) => (
                  <span key={i}>
                    {p}
                    {i < s.AlternateNames.length - 1 && ", "}
                  </span>
                ))}
              </TableCell>
              <TableCell>
                {s.IDsOnRecord.map((p, i) => (
                  <span key={i}>
                    {p}
                    {i < s.IDsOnRecord.length - 1 && ", "}
                  </span>
                ))}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
