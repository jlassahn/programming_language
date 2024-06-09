
#include <stdio.h>

void print_point(double x, double y)
{
	printf("%f,%f\n", x*100.0, 100-y*100.0);
}

typedef struct Vector Vector;
struct Vector
{
	double x;
	double y;
};

typedef struct Matrix Matrix;
struct Matrix
{
	double xx;
	double xy;
	double yx;
	double yy;
};

typedef struct Transform Transform;
struct Transform
{
	Matrix m;
	Vector v;
};

Vector MultMatVec(Matrix m, Vector v)
{
	Vector ret = {
		v.x*m.xx + v.y*m.xy,
		v.x*m.yx + v.y*m.yy
	};
	return ret;
}

Matrix MultMat(Matrix a, Matrix b)
{
	Matrix ret = {
		(a.xx*b.xx + a.xy*b.yx), (a.xx*b.xy + a.xy*b.yy),
		(a.yx*b.xx + a.yy*b.yx), (a.yx*b.xy + a.yy*b.yy)
	};
	return ret;
}

Vector AddVec(Vector a, Vector b)
{
	Vector ret = { a.x + b.x, a.y + b.y };
	return ret;
}

Vector TransformVec(Transform t, Vector v)
{
	return AddVec(t.v, MultMatVec(t.m, v));
}

Transform TransformTrans(Transform a, Transform b)
{
	Transform ret;
	ret.m = MultMat(a.m, b.m);
	ret.v = TransformVec(a, b.v);
	return ret;
}

Transform transforms[] =
{
	{
		.v = { 0, 0 },
		.m = {
			 4.0/16.0, -1.0/16.0,
			 1.0/16.0,  4.0/16.0
		}
	},
	{
		.v = { 4.0/16.0, 1.0/16.0 },
		.m = {
			 1.0/16.0, -4.0/16.0,
			 4.0/16.0,  1.0/16.0
		}
	},
	{
		.v = { 5.0/16.0, 5.0/16.0 },
		.m = {
			 2.0/16.0,  4.0/16.0,
			-4.0/16.0,  2.0/16.0
		}
	},
	{
		.v = { 7.0/16.0, 1.0/16.0 },
		.m = {
			 2.0/16.0,  0.0/16.0,
			 0.0/16.0,  2.0/16.0
		}
	},
	{
		.v = { 9.0/16.0, 1.0/16.0 },
		.m = {
			 2.0/16.0, -4.0/16.0,
			 4.0/16.0,  2.0/16.0
		}
	},
	{
		.v = { 11.0/16.0, 5.0/16.0 },
		.m = {
			 1.0/16.0,  4.0/16.0,
			-4.0/16.0,  1.0/16.0
		}
	},
	{
		.v = { 12.0/16.0, 1.0/16.0 },
		.m = {
			 4.0/16.0,  1.0/16.0,
			-1.0/16.0,  4.0/16.0
		}
	}
};

const int transform_count = sizeof(transforms)/sizeof(Transform);

void do_step(Transform t, int d)
{
	if (d==0)
	{
		Vector p0 = { 0, 0 };
		Vector p1 = { 1, 0 };
		p0 = TransformVec(t, p0);
		p1 = TransformVec(t, p1);
		print_point(p0.x, p0.y);
		print_point(p1.x, p1.y);
	}
	else
	{
		for (int i=0; i<transform_count; i++)
		{
			Transform t_step = TransformTrans(t, transforms[i]);
			do_step(t_step, d-1);
		}
	}
}

int main(void)
{
	puts("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>");
	puts("<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n");

	puts("<svg width=\"100\" height=\"100\" viewBox=\"-10 -10 120 120\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\">\n");

	printf("<polyline points=\"\n");

	Transform t = {
		.v = { 0, 0 },
		.m = {
			1, 0,
			0, 1
		}
	};

	do_step(t, 3);

	printf("\" stroke=\"green\" stroke-width=\"2\" stroke-linejoin=\"round\" fill=\"none\" />\n");

	puts("</svg>");

	return 0;
}

