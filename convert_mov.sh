./main smb.nes encode $1
ffmpeg -f f32le -ar 44100 -ac 1 -i audio.raw audio.mp3
ffmpeg -r 60 -f image2 -i img%06d.png -i audio.mp3 -vcodec libx264 -pix_fmt yuv420p $1.avi
rm audio.raw audio.mp3 img*.png 
