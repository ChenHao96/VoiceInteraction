import javax.sound.sampled.*;
import java.io.*;
import java.util.Arrays;
import java.util.List;

/**
 * 此类线程不安全
 * 此类为录音机的API部分
 * 提供了录音，停止，播放，保存等功能
 */
public abstract class Record4Java implements LineListener {

    public static final int BUF_SIZE = 16 * 1024;

    private boolean play_stop_flag = false;
    private boolean record_stop_flag = false;
    private ByteArrayOutputStream byteArrayOutputStream = null;

    private String recordFilePath;

    public Record4Java() {
    }

    public Record4Java(String recordFilePath) {
        this.recordFilePath = recordFilePath;
    }

    public final byte[] getRecordByteArray() {
        if (byteArrayOutputStream == null) return new byte[0];
        return byteArrayOutputStream.toByteArray();
    }

    //设置AudioFormat的参数
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
        stop,   //调用停止录音的方法
        capture //调用录音的方法
    }

    public final void actionPerformed(Method method) {
        switch (method) {
            case play:
                play(byteArrayOutputStream);
                break;
            case save:
                save(byteArrayOutputStream);
                break;
            case capture:
                byteArrayOutputStream = capture();
                break;
            case stop:
                stop();
                break;
        }
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
            AudioFormat audioFormat = getAudioFormat();
            DataLine.Info info = new DataLine.Info(TargetDataLine.class, audioFormat);
            TargetDataLine targetDataLine = (TargetDataLine) (AudioSystem.getLine(info));
            targetDataLine.open(audioFormat);
            targetDataLine.start();

            record_stop_flag = false;
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            new Thread(new Record(targetDataLine, byteArrayOutputStream)).start();
            return byteArrayOutputStream;
        } catch (LineUnavailableException ex) {
            throw new RuntimeException(ex);
        }
    }

    private void stop() {
        play_stop_flag = true;
        record_stop_flag = true;
    }

    private void play(ByteArrayOutputStream byteArrayOutputStream) {

        if (byteArrayOutputStream == null) return;
        byte audioData[] = byteArrayOutputStream.toByteArray();
        AudioFormat audioFormat = getAudioFormat();
        ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(audioData);
        AudioInputStream audioInputStream = new AudioInputStream(byteArrayInputStream, audioFormat, audioData.length / audioFormat.getFrameSize());

        try {

            DataLine.Info dataLineInfo = new DataLine.Info(SourceDataLine.class, audioFormat);
            SourceDataLine sourceDataLine = (SourceDataLine) AudioSystem.getLine(dataLineInfo);
            sourceDataLine.open(audioFormat);
            sourceDataLine.start();

            play_stop_flag = false;
            new Thread(new Play(audioInputStream, sourceDataLine)).start();
        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            safeClose(byteArrayInputStream);
            safeClose(byteArrayOutputStream);
        }
    }

    private void save(ByteArrayOutputStream byteArrayOutputStream) {

        if (byteArrayOutputStream == null) return;
        byte audioData[] = byteArrayOutputStream.toByteArray();
        AudioFormat audioFormat = getAudioFormat();
        ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(audioData);
        AudioInputStream audioInputStream = new AudioInputStream(byteArrayInputStream, audioFormat, audioData.length / audioFormat.getFrameSize());
        try {

            AudioSystem.write(audioInputStream, AudioFileFormat.Type.WAVE, new File(createSaveFilePath()));
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

    private class Record implements Runnable {

        private TargetDataLine targetDataLine;
        private ByteArrayOutputStream byteArrayOutputStream;

        public Record(TargetDataLine targetDataLine, ByteArrayOutputStream byteArrayOutputStream) {
            this.targetDataLine = targetDataLine;
            this.byteArrayOutputStream = byteArrayOutputStream;
        }

        public void run() {
            targetDataLine.addLineListener(Record4Java.this);
            try {

                byte[] bts = new byte[Record4Java.BUF_SIZE];
                while (true) {
                    int cnt = targetDataLine.read(bts, 0, bts.length);
                    if (cnt > 0) {
                        byteArrayOutputStream.write(bts, 0, cnt);
                    }
                    if (record_stop_flag) {
                        break;
                    }
                }
            } catch (Exception e) {
                throw new RuntimeException(e);
            } finally {
                safeClose(byteArrayOutputStream);
                targetDataLine.removeLineListener(Record4Java.this);
                safeCloseDataLine(targetDataLine);
            }
        }
    }

    private class Play implements Runnable {

        private AudioInputStream audioInputStream;
        private SourceDataLine sourceDataLine;

        public Play(AudioInputStream audioInputStream, SourceDataLine sourceDataLine) {
            this.audioInputStream = audioInputStream;
            this.sourceDataLine = sourceDataLine;
        }

        public void run() {
            try {

                int cnt;
                byte[] bts = new byte[Record4Java.BUF_SIZE];
                while ((cnt = audioInputStream.read(bts)) != -1) {
                    if (cnt > 0) {
                        sourceDataLine.write(bts, 0, cnt);
                    }
                    if (play_stop_flag) {
                        break;
                    }
                }
            } catch (Exception e) {
                throw new RuntimeException(e);
            } finally {
                safeClose(audioInputStream);
                safeCloseDataLine(sourceDataLine);
            }
        }
    }
}
