#include <stdio.h>
#include <opencv2/core/core.hpp>
#include <opencv2/videoio/videoio.hpp>
#include "net.h"
#include "ncnnwrapper.h"

Ncnnnet newNcnnnet() 
{
    return new ncnn::Net(); 
}

void ncnnnetLoad(char* param, char* model, Ncnnnet net)
{
    net->load_param(param);
    net->load_model(model);
}

Rects detectFromByte(unsigned char* data, int cols, int rows, Ncnnnet net, int mode)
{
    const int target_size = 300;

    int img_w = cols;
    int img_h = rows;

    ncnn::Mat in = ncnn::Mat::from_pixels_resize(data, ncnn::Mat::PIXEL_BGR, img_w, img_h, target_size, target_size);

    const float mean_vals[3] = {127.5f, 127.5f, 127.5f};
    const float norm_vals[3] = {1.0/127.5,1.0/127.5,1.0/127.5};
    in.substract_mean_normalize(mean_vals, norm_vals);

    ncnn::Extractor ex = net->create_extractor();
    ex.set_num_threads(2);
    ex.input("data", in);
    ncnn::Mat out;
    ex.extract("detection_out",out);

    std::vector<Rect> vrects;

    Rect rectHeader;
    vrects.push_back(rectHeader);
    for (int i=0; i<out.h; i++)
    {
        const float* values = out.row(i);

        if (values[1] < 0.4) {
            continue;
        }

        if (mode == BODY_DETECT) {
            int label = values[0];
            if (label != 15) continue;
        }

        Rect rect;
        rect.top = values[2] * img_w;
        rect.left = values[3] * img_h;
        rect.width = values[4] * img_w - rect.left;
        rect.height = values[5] * img_h - rect.top;

        vrects.push_back(rect);

        //fprintf(stdout, "x %d y %d x1 %d y1 %d\n", rect.x0, rect.y0, rect.x1, rect.y1);
        //fprintf(stdout, "p: %f\n", values[1]);
    }
    Rects rects = {vrects.size(), &vrects[0]};

    return rects;
}

Rects detectFromPath(Ncnnnet net, int mode, char* camPath) {
    cv::Mat img;
    cv::VideoCapture cap(camPath);
    if(!cap.isOpened()){
        fprintf(stderr, "cap %s failed\n", camPath);
        Rects rects = {0, NULL};
        return rects;
    }  
    cap >> img;
    if(!img.data){
        fprintf(stderr, "img %s failed\n", camPath);
        Rects rects = {0, NULL};
        return rects;
    }

    Rects rects = detectFromByte(img.data, img.cols, img.rows, net, mode);

    img.release();
    cap.release();
    return rects;
}

