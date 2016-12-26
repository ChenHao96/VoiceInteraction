import javax.sound.sampled.*;
import java.io.*;

/**
 * 此类为录音机的API部分
 * 提供了录音，停止，播放，保存等功能
 */
public abstract class Record4Java implements RecordDataListener {

    /**
     * 指定的文件保存路径
     */
    private String recordFilePath;
    private boolean play_stop_flag = false;
    private boolean record_stop_flag = false;
    private ByteArrayOutputStream byteArrayOutputStream = null;

    public Record4Java() {
    }

    public Record4Java(String recordFilePath) {
        this.recordFilePath = recordFilePath;
    }

    public String getRecordFilePath() {
        return recordFilePath;
    }

    public void setRecordFilePath(String recordFilePath) {
        this.recordFilePath = recordFilePath;
    }

    public final boolean isRecordStop() {
        return record_stop_flag;
    }

    public final boolean isPlayStop() {
        return play_stop_flag;
    }
    
    /**
     * 设置AudioFormat的参数
     */
    public AudioFormat getAudioFormat() {

        //采样率是每秒播放和录制的样本数
        // 采样率8000,11025,16000,22050,44100
        float sampleRate = 16000.0F;
        //sampleSizeInBits表示每个具有此格式的声音样本中的位数
        int sampleSizeInBits = 16;// 8,16
        int channels = 1;// 单声道为1，立体声为2
        boolean signed = true;// true:Encoding.PCM_SIGNED,false:Encoding.PCM_UNSIGNED
        boolean bigEndian = true;

        return new AudioFormat(sampleRate, sampleSizeInBits, channels, signed, bigEndian);
    }

    public enum Method {
        play,   //调用播放录音的方法
        save,   //调用保存录音的方法
        stop,   //调用停止录音、播放的方法
        capture //调用录音的方法
    }

    /**
     * 录音机的录音，停止，播放，保存等功能
     *
     * @param method 执行的方法
     * @param waves  播放的wav文件(只支持wav)
     * @return 调用save方法时返回文件路径
     */
    public final String actionPerformed(Method method, File... waves) {
        switch (method) {
            case play:
                play(byteArrayOutputStream, waves);
                break;
            case save:
                return save(byteArrayOutputStream);
            case capture:
                byteArrayOutputStream = capture();
                break;
            case stop:
                stop();
                break;
        }

        return null;
    }

    public static void safeClose(Closeable closeable) {
        if (closeable != null) {
            try {
                closeable.close();
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
    }

    private static void safeCloseDataLine(DataLine dataLine) {
        if (dataLine != null) {
            try {
                dataLine.close();
                dataLine.drain();
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        }
    }

    private ByteArrayOutputStream capture() {

        try {
            record_stop_flag = false;
            TargetDataLine targetDataLine = getTargetDataLine();
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            try {

                byte[] bts = new byte[BUF_SIZE];
                do {
                    int cnt = targetDataLine.read(bts, 0, bts.length);
                    if (cnt > 0) {
                        byteArrayOutputStream.write(bts, 0, cnt);
                    }
                    processRecordData(bts);
                } while (!record_stop_flag);
            } finally {
                safeClose(byteArrayOutputStream);
                safeCloseDataLine(targetDataLine);
            }

            return byteArrayOutputStream;
        } catch (Exception ex) {
            throw new RuntimeException(ex);
        }
    }

    private TargetDataLine getTargetDataLine() throws LineUnavailableException {

        AudioFormat audioFormat = getAudioFormat();
        DataLine.Info info = new DataLine.Info(TargetDataLine.class, audioFormat);
        TargetDataLine targetDataLine = (TargetDataLine) (AudioSystem.getLine(info));
        targetDataLine.open(audioFormat);
        targetDataLine.start();

        return targetDataLine;
    }

    private void stop() {
        play_stop_flag = true;
        record_stop_flag = true;
    }

    private void play(ByteArrayOutputStream byteArrayOutputStream, File... waves) {

        if (byteArrayOutputStream != null) {

            byte audioData[] = byteArrayOutputStream.toByteArray();
            AudioFormat audioFormat = getAudioFormat();
            ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(audioData);
            AudioInputStream audioInputStream = new AudioInputStream(byteArrayInputStream, audioFormat, audioData.length / audioFormat.getFrameSize());
            try {

                playWave(audioFormat, audioInputStream);
            } catch (Exception e) {
                throw new RuntimeException(e);
            } finally {
                safeClose(byteArrayInputStream);
                safeClose(byteArrayOutputStream);
            }
        }

        if (waves != null && waves.length > 0) {

            try {
                for (int i = 0; i < waves.length; i++) {

                    AudioInputStream audioInputStream = AudioSystem.getAudioInputStream(waves[i]);
                    playWave(audioInputStream.getFormat(), audioInputStream);
                }
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        }
    }

    private void playWave(AudioFormat audioFormat, AudioInputStream audioInputStream) throws LineUnavailableException, IOException {

        DataLine.Info dataLineInfo = new DataLine.Info(SourceDataLine.class, audioFormat);
        SourceDataLine sourceDataLine = (SourceDataLine) AudioSystem.getLine(dataLineInfo);
        sourceDataLine.open(audioFormat);
        sourceDataLine.start();

        try {
            int cnt;
            play_stop_flag = false;
            byte[] bts = new byte[BUF_SIZE];
            while ((cnt = audioInputStream.read(bts)) != -1) {
                if (cnt > 0) {
                    sourceDataLine.write(bts, 0, cnt);
                }
                if (play_stop_flag) {
                    break;
                }
            }
            play_stop_flag = true;
        } finally {
            safeClose(audioInputStream);
            safeCloseDataLine(sourceDataLine);
        }
    }

    private String save(ByteArrayOutputStream byteArrayOutputStream) {

        if (byteArrayOutputStream == null) return null;
        byte audioData[] = byteArrayOutputStream.toByteArray();
        AudioFormat audioFormat = getAudioFormat();
        ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(audioData);
        AudioInputStream audioInputStream = new AudioInputStream(byteArrayInputStream, audioFormat, audioData.length / audioFormat.getFrameSize());
        try {

            String filePath = createSaveFilePath();
            AudioSystem.write(audioInputStream, AudioFileFormat.Type.WAVE, new File(filePath));
            return filePath;
        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            safeClose(audioInputStream);
            safeClose(byteArrayInputStream);
            safeClose(byteArrayOutputStream);
        }
    }

    private String createSaveFilePath() {

        String path = System.getProperty("user.home");
        if (recordFilePath != null && recordFilePath.trim().length() > 0) {
            path = recordFilePath;
            if (path.lastIndexOf("/") == path.length() - 1 || path.lastIndexOf("\\") == path.length() - 1) {
                path = path.substring(0, path.length() - 1);
            }
        }

        File filePath = new File(path + "/record4java");
        if (!filePath.exists()) {
            if (filePath.mkdirs()) {
                System.out.println("文件路径不存在，已成功创建! 文件存储路径:" + filePath.getPath());
            }
        }

        return filePath.getPath() + "/" + System.currentTimeMillis() + ".wav";
    }
}
