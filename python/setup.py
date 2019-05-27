from setuptools import setup, Distribution

class BinaryDistribution(Distribution):
    def has_ext_modules(foo):
        return True

setup(
    name='visenze_recordio',
    version='0.1',
    description='An implementation of the RecordIO file format.',
    url='https://github.com/PaddlePaddle/recordio',
    author='PaddlePaddle Authors',
    author_email='paddle-dev@baidu.com',
    license='Apache 2.0',
    packages=['recordio'],
    package_data={
        'recordio': ['librecordio.so'],
    },
    distclass=BinaryDistribution
)
