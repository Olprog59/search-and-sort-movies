# Search And Sort Movies or Series

This application will allow you to sort films and series.
It will clean up the names and move them to the folder you want.

Ex: 
- /be_sorted/movie_sam_2020_to$http://sAm.EN-01.mkv => /movies/movie-sam-2020.mkv
- /be_sorted/serie_S1_e12_qWerTy_aZerty.mKv => /series/fringe/season-1/fringe-S01-E12.mkv

### Choose your Volumes :
- /be_sorted
- /movies
- /series

Ex.:
```bash
docker container run \
-v /mnt/user/dlna/be_sorted:/be_sorted \
-v /mnt/user/dlna/movies:/movies \
-v /mnt/user/dlna/series:/series \
--name movies media-organizer
```
