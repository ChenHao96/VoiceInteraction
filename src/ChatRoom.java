import javax.swing.*;
import javax.swing.border.TitledBorder;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.IOException;
import java.net.*;

public class ChatRoom {

    private DatagramSocket socket = null;
    private DatagramSocket server = null;
    private JTextField userName;
    private JButton linkButton;
    private JButton closeButton;
    private JTextArea area;
    private JTextField sendFile;
    private JButton sendButton;
    private ThreadRead read = null;

    public ChatRoom() {

        /**
         * 创建窗口
         */
        JFrame frame = new JFrame("多人聊天室");
        frame.setSize(500, 500);//设置窗口的大小
        frame.setResizable(false);//不可改变窗口大小
        //noinspection MagicConstant
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);//设置窗口的关闭方式
        frame.setLocationRelativeTo(null);//设置窗口的位置在屏幕的中间(居中)

        /**
         * 用户配置框
         */
        JLabel label = new JLabel("用户名:");//静态标签  显示"用户名:"
        userName = new JTextField(10);//单行文本框  获取用户名
        linkButton = new JButton("连接");//创建按钮  连接按钮
        closeButton = new JButton("退出");//创建按钮  退出按钮
        closeButton.setEnabled(false);//设置退出按钮不可用  未连接到聊天室不能执行退出
        JPanel panel = new JPanel();
        panel.add(label);
        panel.add(userName);
        panel.add(linkButton);
        panel.add(closeButton);
        panel.setBorder(new TitledBorder("用户配置"));//设置标题边框

        /**
         * 信息框
         */
        area = new JTextArea();//多行文本域  显示接收到的消息
        area.setLineWrap(true);//设置自动换行
        area.setEditable(false);//设置不可输入  只用于显示内容
        JScrollPane pane = new JScrollPane(area);//添加滚动条
        pane.setBorder(new TitledBorder("聊天消息"));//设置标题边框

        /**
         * 发送框
         */
        sendFile = new JTextField(30);//单行文本框  获取要发送的内容
        sendButton = new JButton("发送");//创建按钮  发送按钮
        sendButton.setEnabled(false);//设置按钮不可用   未连接到聊天室不能执行发送
        JPanel panel2 = new JPanel();
        panel2.add(sendFile);
        panel2.add(sendButton);
        panel2.setBorder(new TitledBorder("消息输入区"));//设置标题边框

        ListenerAll listener = new ListenerAll();
        userName.addActionListener(listener);
        linkButton.addActionListener(listener);
        sendButton.addActionListener(listener);
        sendFile.addActionListener(listener);
        closeButton.addActionListener(listener);

        /**
         * 把组件添加到窗口中
         * 设置显示窗口
         */
        frame.add(panel, BorderLayout.NORTH);
        frame.add(pane, BorderLayout.CENTER);
        frame.add(panel2, BorderLayout.SOUTH);
        frame.setVisible(true);
    }

    public class ListenerAll implements ActionListener {

        private static final int PORT = 60000;
        private final String HOSTNAME;

        public ListenerAll() {
            InetAddress addr;
            try {
                addr = InetAddress.getLocalHost();
            } catch (UnknownHostException e) {
                throw new RuntimeException(e);
            }
            String ip = addr.getHostAddress();//获得本机IP
            HOSTNAME = ip.substring(0, ip.lastIndexOf(".") + 1) + "255";
        }

        @SuppressWarnings("deprecation")
        @Override
        public void actionPerformed(ActionEvent e) {

            if (e.getSource() == linkButton || e.getSource() == userName) {

                if (userName.getText().trim().isEmpty()) {
                    area.append(userName.getText().trim() + "请输入用户名！\r\n");
                    return;
                }
                try {
                    socket = new DatagramSocket();
                    byte[] buf = (userName.getText().trim() + " 上线了\r\n").getBytes();
                    socket.send(new DatagramPacket(buf, 0, buf.length, new InetSocketAddress(HOSTNAME, PORT)));
                } catch (IOException e2) {
                    e2.printStackTrace();
                }

                //启动线程
                try {

                    server = new DatagramSocket(PORT);
                    read = new ThreadRead(area, server);
                    read.start();
                } catch (SocketException e3) {
                    e3.printStackTrace();
                }

                //释放退出、发送按钮
                closeButton.setEnabled(true);
                sendButton.setEnabled(true);
                userName.setEnabled(false);
                linkButton.setEnabled(false);

            }
            if (e.getSource() == sendButton || e.getSource() == sendFile) {
                try {

                    socket = new DatagramSocket();
                    if (sendFile.getText().trim().isEmpty()) {
                        area.append("输入不能为空！！！\r\n");
                        return;
                    }
                    byte[] buf = (userName.getText().trim() + ":" + sendFile.getText().trim()).getBytes();
                    sendFile.setText(null);
                    socket.send(new DatagramPacket(buf, 0, buf.length, new InetSocketAddress(HOSTNAME, PORT)));
                } catch (IOException e3) {
                    e3.printStackTrace();
                }
            }
            if (e.getSource() == closeButton) {

                if (read != null) {
                    read.stop();
                }
                if (socket != null) {
                    socket.close();
                }
                if (server != null) {
                    server.close();
                }
                if (socket != null && server != null && read != null) {

                    closeButton.setEnabled(false);
                    sendButton.setEnabled(false);
                    try {

                        socket = new DatagramSocket();
                        byte[] buf = (userName.getText().trim() + " 退出聊天室\r\n").getBytes();
                        socket.send(new DatagramPacket(buf, 0, buf.length, new InetSocketAddress(HOSTNAME, PORT)));
                    } catch (IOException e4) {
                        e4.printStackTrace();
                    }
                    linkButton.setEnabled(true);
                    userName.setText(null);
                    userName.setEnabled(true);
                    area.setText(null);
                }
            }
        }
    }

    public class ThreadRead extends Thread {

        private JTextArea area;

        private DatagramSocket server;

        public ThreadRead(JTextArea area, DatagramSocket server) {

            this.area = area;
            this.server = server;
        }

        @Override
        public void run() {
            try {
                //noinspection InfiniteLoopStatement
                for (; ; ) {
                    byte[] buf = new byte[1024];
                    server.receive(new DatagramPacket(buf, buf.length));
                    area.append(new String(buf).trim() + "\r\n");
                }
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
}
