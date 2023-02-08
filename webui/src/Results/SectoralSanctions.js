import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Typography from '@material-ui/core/Typography';
import Paper from "@material-ui/core/Paper";

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
