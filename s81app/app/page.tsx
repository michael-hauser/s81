import React from 'react';
import ThemeButton from './components/ThemeButton/ThemeButton';
import styles from './page.module.scss';
import Main from './main';

const Home: React.FC = () => {
  return (
    <div className={styles.app}>
      <header className={styles.header}>
        <div>
          <h1>81st Street-Museum of<br />Natural History station</h1>
          <div className=''>Tuesday, July 30 | 8:46 AM</div>
        </div>
        <ThemeButton />
      </header>
      <Main />
    </div>
  );
};

export default Home;
