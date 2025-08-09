/**
 *    SPDX-License-Identifier: Apache-2.0
 */

import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import { Box, Typography, Container, Link } from '@material-ui/core';
import SchoolIcon from '@material-ui/icons/School';
import FacebookIcon from '@material-ui/icons/Facebook';
import RoomIcon from '@material-ui/icons/Room';
import GitHubIcon from '@material-ui/icons/GitHub';
import Logo from '../../static/images/logoKMA.png';

const styles = theme => {
  const { type } = theme.palette;
  const dark = type === 'dark';
  return {
    root: {
      width: '100%',
      borderTop: dark ? '1px solid #333' : '1px solid #e0e0e0',
      backgroundColor: dark ? theme.palette.background.default : '#f5f5f5',
      paddingTop: theme.spacing(4),
      paddingBottom: theme.spacing(4),
	  marginTop: 16
    },
    container: {
      display: 'flex',
      flexDirection: 'column',
      gap: theme.spacing(2),
    },
    logoContainer: {
      display: 'flex',
      alignItems: 'center',
      gap: theme.spacing(1),
    },
    logo: {
      height: 32,
      width: 32,
    },
    logoText: {
      fontWeight: 600,
      color: 'red',
    },
    description: {
      fontSize: '0.875rem',
      color: dark ? theme.palette.text.secondary : '#6c757d',
    },
    footerBottom: {
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'space-between',
      gap: theme.spacing(2),
      [theme.breakpoints.up('md')]: {
        flexDirection: 'row',
      },
    },
    copyright: {
      fontSize: '0.875rem',
      color: dark ? theme.palette.text.secondary : '#6c757d',
    },
    socialLinks: {
      display: 'flex',
      alignItems: 'center',
      gap: theme.spacing(3),
    },
    socialIcon: {
      color: dark ? theme.palette.text.secondary : '#6c757d',
      '&:hover': {
        color: theme.palette.primary.main,
      },
    },
    versionInfo: {
      backgroundColor: dark ? '#5e558e' : '#e8e8e8',
      color: dark ? '#ffffff' : undefined,
      textAlign: 'center',
      position: 'fixed',
      left: 0,
      right: 0,
      bottom: 0,
      padding: theme.spacing(0.5),
    },
  };
};

const socialLinks = [
  { icon: SchoolIcon, href: 'https://actvn.edu.vn', label: 'Trang chủ học viện' },
  { icon: FacebookIcon, href: 'https://www.facebook.com/hocvienkythuatmatma', label: 'Facebook' },
  { icon: RoomIcon, href: 'https://maps.app.goo.gl/nH4ungjtTKWfox2c8', label: 'Địa chỉ' },
  { icon: GitHubIcon, href: 'https://github.com/vnkmasc/Kmasc', label: 'Github' }
];

const FooterView = ({ classes }) => (
  <Box className={classes.root}>
    <Container className={classes.container}>
      <Box className={classes.logoContainer}>
        <Link href="/">
          <img src={Logo} alt="logo" title="Kmasc" className={classes.logo} />
        </Link>
        <Typography variant="h6" className={classes.logoText}>
          Kmasc
        </Typography>
      </Box>
      
      <Typography className={classes.description}>
        Giải pháp quản lý văn bằng chứng chỉ ứng dụng Blockchain.
      </Typography>
      
      <Box className={classes.footerBottom}>
        <Typography className={classes.copyright}>
          © 2025 Kmasc. Bản quyền thuộc về khoa CNTT Học Viện Kỹ Thuật Mật Mã, phát triển dựa trên Hyperledger Explorer.
        </Typography>
        
        <Box className={classes.socialLinks}>
          {socialLinks?.map((social, idx) => (
            <Link 
              key={idx} 
              href={social.href} 
              aria-label={social.label} 
              target="_blank" 
              rel="noopener"
            >
              <social.icon className={classes.socialIcon} />
            </Link>
          ))}
        </Box>
      </Box>
    </Container>
  </Box>
);

export default withStyles(styles)(FooterView);
